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

type ImageDB struct{}

func (c *ImageDB) Upload(file io.Reader, dataType string, uid primitive.ObjectID) (primitive.ObjectID, error) {
	bucket, err := gridfs.NewBucket(
		Instance.mongoClient.Database("Images"),
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

func (c *ImageDB) DownloadById(id primitive.ObjectID) (bytes.Buffer, error) {
	bucket, _ := gridfs.NewBucket(
		Instance.mongoClient.Database("Images"),
	)
	buf := bytes.Buffer{}
	_, err := bucket.DownloadToStream(id, &buf)
	return buf, err
}

func (c *ImageDB) GetMetdataById(id primitive.ObjectID) (map[string]interface{}, error) {
	fsCollections := Instance.mongoClient.Database("Images").Collection("fs.files")
	metadata := make(map[string]interface{})
	if err := fsCollections.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&metadata); err != nil {
		return metadata, err
	}
	return metadata["metadata"].(map[string]interface{}), nil
}

func (c *ImageDB) DeleteImageById(id primitive.ObjectID) error {
	bucket, _ := gridfs.NewBucket(
		Instance.mongoClient.Database("Images"),
	)
	return bucket.Delete(id)
}

func (c *ImageDB) ImageList(uploaderId primitive.ObjectID) ([]bson.M, error) {
	fsCollections := Instance.mongoClient.Database("Images").Collection("fs.files")
	var data []bson.M
	matchStage := bson.D{{"$match", bson.D{{"metadata.uploader", uploaderId}}}}
	mapStage := bson.D{{"$project", bson.D{{"id", "$_id"}, {"_id", 0}, {"uploadDate", "$uploadDate"}}}}

	cursor, err := fsCollections.Aggregate(timeoutContext(),
		mongo.Pipeline{matchStage, mapStage})
	if err != nil {
		return data, err
	}
	err = cursor.All(timeoutContext(), &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
