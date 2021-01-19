package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func parseMatchesToIdArray(bodyMap bson.M) []string {
	res := []string{}
	for _, match := range bodyMap["matches"].([]interface{}) {
		gameRawId := match.(map[string]interface{})["gameId"]
		id := strconv.Itoa(int(gameRawId.(float64)))
		res = append(res, id)
	}
	return res
}

func increaseMatchMap(mutex *sync.Mutex, wg *sync.WaitGroup, countMap *map[string]int32, accountId string, timestamp int64) {
	for i := 1; i <= 3600; i++ {
		defer (*wg).Done()
		respMap, respCode := utils.RequestToRiotServer("/lol/match/v4/matchlists/by-account/"+accountId,
			bson.M{"beginTime": strconv.FormatInt(timestamp, 10)})
		if respCode == 429 {
			log.Print("RateLimit exceded")
			log.Print("increaseMatchMap() Retry: " + strconv.Itoa(i))
			time.Sleep(3 * time.Second)
			continue
		} else if respCode != 200 {
			log.Print("RIOT SERVER RESPONSE AS " + strconv.Itoa(respCode))
			return
		}
		array := parseMatchesToIdArray(respMap)
		log.Print(array)
		for _, id := range array {
			(*mutex).Lock()
			(*countMap)[id] += 1
			(*mutex).Unlock()
		}
		return
	}
}

func requestAndStoreToDB(mutex *sync.Mutex, wg *sync.WaitGroup, gameId string, userExists map[string]bool, results *[]primitive.ObjectID) {
	defer (*wg).Done()
	for i := 1; i <= 3600; i++ {
		respMap, respCode := utils.RequestToRiotServer("/lol/match/v4/matches/"+gameId, nil)
		if respCode == 429 {
			log.Print("RateLimit exceded")
			log.Print(respMap)
			log.Print("Retry: " + strconv.Itoa(i))
			time.Sleep(3 * time.Second)
			continue
		} else if respCode != 200 {
			log.Print("RIOT RESPONSE " + strconv.Itoa(respCode))
			log.Print(respMap)
			return
		}

		var participateId int
		var win bool
		for _, v := range respMap["participantIdentities"].([]interface{}) {
			vMap := v.(map[string]interface{})
			participaint := vMap["player"].(map[string]interface{})
			accountId := participaint["accountId"].(string)
			if userExists[accountId] {
				participateId = int(vMap["participantId"].(float64))
				break
			}
		}
		for _, v := range respMap["participants"].([]interface{}) {
			vMap := v.(map[string]interface{})
			id := int(vMap["participantId"].(float64))
			if id == participateId {
				stats := vMap["stats"].(map[string]interface{})
				win = stats["win"].(bool)
			}
		}
		timestamp := int64(respMap["gameCreation"].(float64))
		mutex.Lock()
		id, _ := database.Instance.LOLHistory.AddLolHistory(respMap, win, timestamp)
		(*results) = append((*results), id)
		mutex.Unlock()
		return
	}
}

// @Summary Update lol history
// @Description Find game after the timestamp and store to DB
// @Description **Application Token** is required.
// @Accept  json
// @Produce  json
// @Param X-Dfd-App-Auth header string true "application access token"
// @Success 200 {object} docmodels.ResponseSuccess{results=[]primitive.ObjectID} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/lol/history
// @Router /api/v1/lol/history/updater [post]
func PostLolHistoryUpdate(c *gin.Context) {
	_, respCode := utils.RequestToRiotServer("/lol/status/v4/platform-data", nil)

	if respCode == 403 {
		log.Print("RIOT API Token is invalid")
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}

	timeString, err := database.Instance.Redis.Get("UpdateTimestamp")
	if err != nil {
		timeString = time.Now().Format(time.RFC3339)
	}

	_, err = time.Parse(time.RFC3339, timeString)
	if err != nil {
		log.Print(timeString)
		log.Fatal(err.Error())
		return
	}
	database.Instance.Redis.Set("UpdateTimestamp", time.Now().Format(time.RFC3339))

	users := database.Instance.User.GetLoLInfoList()
	userExistsMap := map[string]bool{}
	countMap := map[string]int32{}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for _, user := range users {
		wg.Add(1)
		userExistsMap[user["lol_account_id"].(string)] = true
		go increaseMatchMap(&mutex, &wg, &countMap, user["lol_account_id"].(string), 1610285111000)
	}
	wg.Wait()
	wg = sync.WaitGroup{}
	mutex = sync.Mutex{}
	results := []primitive.ObjectID{}
	for k, v := range countMap {
		if v >= 3 {
			wg.Add(1)
			go requestAndStoreToDB(&mutex, &wg, k, userExistsMap, &results)
		}
	}
	wg.Wait()
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"results": results}))
	return
}

// @Summary Update lol history
// @Description Find game after the timestamp and store to DB
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{games=[]docmodels.ResponseLoLHistory} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @tags api/v1/lol/history
// @Router /api/v1/lol/historys [get]
func GetLolHistoryList(c *gin.Context) {
	games := database.Instance.LOLHistory.GetList()
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"games": games}))
	return
}

// @Summary Update lol history
// @Description Find game after the timestamp and store to DB
// @Accept  json
// @Produce  json
// @Param id path string true "Game ID"
// @Success 200 {object} docmodels.ResponseSuccess{game=interface{}} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @tags api/v1/lol/history
// @Router /api/v1/lol/history/{id} [get]
func GetLolHistory(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	history, err := database.Instance.LOLHistory.GetLolHistory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("game not found"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"game": history}))
	return
}
