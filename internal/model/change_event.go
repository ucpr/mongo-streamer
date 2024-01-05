package model

type (
	// ChangeEvent is a struct that represents a change stream event.
	ChangeEvent struct {
		ID                string             `avro:"_id" bson:"_id" json:"_id"`
		OperationType     string             `avro:"operationType" bson:"operation_type" json:"operation_type"`
		FullDocument      []byte             `avro:"fullDocument" bson:"full_document" json:"full_document"`
		DocumentKey       string             `avro:"documentKey" bson:"document_key" json:"document_key"`
		UpdateDescription *UpdateDescription `avro:"updateDescription" bson:"update_description" json:"update_description"`
		Namespace         Namespace          `avro:"ns" bson:"namespace" json:"namespace"`
		To                *Namespace         `avro:"to" bson:"to" json:"to"`
	}

	// UpdateDescription is a struct that represents an update description of change stream event.
	UpdateDescription struct {
		UpdatedFields string `avro:"updatedFields" bson:"updated_fields" json:"updated_fields"`
		RemovedFields string `avro:"removedFields" bson:"removed_fields" json:"removed_fields"`
	}

	// Namespace is a struct that represents a namespace of change stream event.
	Namespace struct {
		DB   string `avro:"db" bson:"db" json:"db"`
		Coll string `avro:"coll" bson:"coll" json:"coll"`
	}
)
