package database

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
)

//ImageDB is a structure for image data access
type ImageDB struct {
	BaseDB
	database *mongo.Database
}

//Upload image into GridFS database
func (db *ImageDB) Upload(file io.Reader, dataType string, uid primitive.ObjectID) (primitive.ObjectID, error) {
	bucket, err := gridfs.NewBucket(
		db.database,
	)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	meta := make(map[string]interface{})
	meta["uploader"] = uid
	meta["content-type"] = dataType
	return bucket.UploadFromStream(
		"", file, &options.UploadOptions{
			Metadata: meta,
		},
	)

}

//DownloadByID image from GridFS database
func (db *ImageDB) DownloadByID(id primitive.ObjectID) (bytes.Buffer, error) {
	bucket, _ := gridfs.NewBucket(
		db.database,
	)
	buf := bytes.Buffer{}
	_, err := bucket.DownloadToStream(id, &buf)
	return buf, err
}

//GetMetdataByID from GridFS database
func (db *ImageDB) GetMetdataByID(id primitive.ObjectID) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})
	if err := db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&metadata); err != nil {
		return metadata, err
	}
	return metadata["metadata"].(map[string]interface{}), nil
}

//DeleteImageByID from GridFS database
func (db *ImageDB) DeleteImageByID(id primitive.ObjectID) error {
	bucket, _ := gridfs.NewBucket(
		db.database,
	)
	return bucket.Delete(id)
}

//ImageList get image id list uploaded by uploaderID from GridFS database
func (db *ImageDB) ImageList(uploaderID primitive.ObjectID) ([]bson.M, error) {
	var data []bson.M
	matchStage := bson.D{{"$match", bson.D{{"metadata.uploader", uploaderID}}}}
	projectStage := bson.D{{"$project", bson.D{{"id", "$_id"}, {"_id", 0}, {"uploadDate", "$uploadDate"}}}}

	cursor, err := db.collection.Aggregate(timeoutContext(),
		mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		return data, err
	}
	err = cursor.All(timeoutContext(), &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
