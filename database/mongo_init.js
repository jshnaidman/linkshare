db_name = DB_NAME
schema_version = SCHEMA_VERSION
session_lifetime_seconds = SESSION_LIFETIME_SECONDS

db.createCollection('pages', {
    validator: {
        $jsonSchema: {
            required: ["schema", "dateAdded", "user_id"],
            properties: {
                _id: {
                    bsonType: "string",
                    maxLength: 30,
                    pattern: "[A-Za-z0-9\\-_]",
                    description: "A URL for the page. Maximum 30 characters"
                },
                dateAdded: {
                    bsonType: "date",
                    description: "The time at which the page was added."
                },
                description: {
                    bsonType: "string"
                },
                title: {
                    bsonType: "string"
                },
                links: {
                    bsonType: "array",
                    items: {
                        bsonType: "string",
                        maxLength: 2048,
                        description: "A URL. The max length of a URL is 2048 characters"
                    },
                    maxItems: 200,
                },
                user_id: {
                    bsonType: "objectId",
                    description: "The user_id of whoever owns the page"
                },
                schema: {
                    bsonType: "int",
                    description: "Schema version"
                }
            }
        }
    }
})

db.createCollection('users', {
    validator: {
        $jsonSchema: {
            required: ["schema"],
            properties: {
                _id: {
                    bsonType: "objectId",
                    description: "user ID"
                },
                username: {
                    bsonType: "string",
                    description: "Username. Can change, so can't be _id"
                },
                // google_id is not indexed because we only need it on login
                google_id: {
                    bsonType: "string",
                    description: "The userID for google login"
                },
                email: {
                    bsonType: "string"
                },
                first_name: {
                    bsonType: "string"
                },
                last_name: {
                    bsonType: "string"
                },
                pages: {
                    bsonType: "array",
                    items: {
                        bsonType: "string",
                        description: "The pages owned by the user."
                    },
                    maxItems: 200,
                },
                schema: {
                    bsonType: "int",
                    description: "Schema version"
                }
            }
        }
    }
})

// usernames queried often
db.users.createIndex({ username: 1 }, {
    unique: true, partialFilterExpression: {
        username: { $type: "string" }
    }
})

db.createCollection('sessions', {
    validator: {
        $jsonSchema: {
            required: ["schema", "user_id"],
            properties: {
                _id: {
                    bsonType: "string",
                    description: "session ID"
                },
                user_id: {
                    bsonType: "objectId",
                    description: "The associated user_id of the session"
                },
                modified: {
                    bsonType: "date",
                    description: "Last date modified"
                },
                schema: {
                    bsonType: "int",
                    description: "Schema version"
                },
            }
        }
    }
})

// For now I'm just going to set session to expire after 24 hours, might implement
// something better later. TODO 
db.sessions.createIndex({ modified: 1 }, { sparse: true, expireAfterSeconds: session_lifetime_seconds })