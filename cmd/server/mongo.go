package main

import (
	"context"
	"log"
	"time"

	pb "github.com/nazufel/wizter/wizard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// storage struct holds information about connecting to storage
type storage struct {
	c  *mongo.Client
	db *mongo.Database
}

// func dbConnect creates a new connection to storage
func dbConnect() (*storage, error) {

	s := new(storage)

	dbConnectionString := "mongodb://mongo.default.svc.cluster.local"

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnectionString))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	// test connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error pinging DB: %v", err)
	}

	s.c = client
	s.db = client.Database("wizards")

	return s, nil
}

func (s *storage) getWizards() ([]pb.Wizard, error) {
	// set find options behavior
	findOptions := options.Find()
	// findOptions.SetLimit(25)

	// filter by company, this is the default behavior
	filter := bson.D{{}}

	cursor, err := s.db.Collection("wizards").Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Printf("failed to get wizards: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var wizards []pb.Wizard

	for cursor.Next(context.Background()) {
		var wizard pb.Wizard

		err := cursor.Decode(&wizard)
		if err != nil {
			log.Printf("unable to decode wizard cursor into struct: %v", err)
			return nil, err
		}

		wizards = append(wizards, wizard)

		// comment this log statement as part of the server demo
		log.Printf("sending wizard to client: %v", wizard.GetName())

	}

	err = cursor.Err()
	if err != nil {
		log.Printf("error with the client cursor: %v", err)
		return nil, err
	}

	return wizards, nil
}

// seedData drops the wizards collection and seeds it with fresh data to the demo
func (s *storage) seedData() error {

	// // drop the collection in order to see fresh data for a new run
	err := s.db.Collection("wizards").Drop(context.Background())
	if err != nil {
		log.Fatalf("unable to drop the wizard collection %v", err)
	}

	// seed the database

	wizards := []pb.Wizard{
		{Name: "Harry Potter",
			House:      "Gryffindor",
			DeathEater: false,
		},
		{Name: "Ron Weasley",
			House:      "Gryffindor",
			DeathEater: false,
		},
		{Name: "Hermione Granger",
			House:      "Gryffindor",
			DeathEater: false,
		},
		{Name: "Cho Chang",
			House:      "Ravenclaw",
			DeathEater: false,
		},
		{Name: "Luna Lovegood",
			House:      "Ravenclaw",
			DeathEater: false,
		},
		{Name: "Sybill Trelawney",
			House:      "Ravenclaw",
			DeathEater: false,
		},
		{Name: "Pomona Sprout",
			House:      "Hufflepuff",
			DeathEater: false,
		},
		{Name: "Cedric Diggory",
			House:      "Hufflepuff",
			DeathEater: false,
		},
		{Name: "Newton Scamander",
			House:      "Hufflepuff",
			DeathEater: false,
		},
		{Name: "Draco Malfoy",
			House:      "Slytherin",
			DeathEater: true,
		},
		{Name: "Bellatrix Lestrange",
			House:      "Slytherin",
			DeathEater: true,
		},
		{Name: "Severus Snape",
			House:      "Slytherin",
			DeathEater: false,
		},
	}
	log.Printf("connected to the database")

	for i := range wizards {
		_, err = s.db.Collection("wizards").InsertOne(context.Background(), wizards[i])
		if err != nil {
			log.Fatalf("failed to insert document: %v", err)
		}
		log.Printf("inserted document for wizard: %s", wizards[i].Name)
		i++
	}

	log.Printf("inserted %v documents", len(wizards))

	return err
}