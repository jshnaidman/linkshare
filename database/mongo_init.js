db_name = DB_NAME
schema_version = SCHEMA_VERSION
num_urls = URL_CAPACITY
bulk_load_count = RUN_ONCE_BULK_LOAD_PAGE_AMOUNT

db.createCollection('pages', {
    validator: {
        $jsonSchema: {
            required: ["schema"],
            properties: {
                _id: {
                    bsonType: "int",
                    description: "pages which are used"
                },
                alias: {
                    bsonType: "string",
                    maxLength: 30,
                    pattern: "[A-Za-z0-9\\-_]",
                    description: "A custom URL for the page. Maximum 30 characters"
                },
                dateAdded: {
                    bsonType: "date",
                    description: "Date added. If missing, page is free"
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
                    bsonType: "int",
                    description: "The user_id of whoever owns the page"
                },
                schema: {
                    bsonType: "int"
                }
            }
        }
    }
})

db.createCollection('unusedPagesIDs', {
    validator: {
        $jsonSchema: {
            required: ["schema"],
            properties: {
                _id: {
                    bsonType: "int",
                    description: "Unused Page IDs"
                },
                schema: {
                    bsonType: "int"
                }
            }
        }
    }
})

db.createCollection('users', {
    validator: {
        $jsonSchema: {
            required: ["username", "schema"],
            properties: {
                _id: {
                    bsonType: "int",
                    description: "Unchangeable user ID"
                },
                username: {
                    bsonType: "string",
                    description: "Chosen username"
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
                        bsonType: "int",
                        description: "The pages owned by the user."
                    },
                    maxItems: 200,
                },
                schema: {
                    bsonType: "int"
                }
            }
        }
    }
})

// usernames must be unique
db.users.createIndex({ username: 1 }, { unique: true })

// aliases must be unique, but we should allow multiple pages to have null aliases
db.pages.createIndex({ alias: 1 }, {
    unique: true, partialFilterExpression: {
        alias: { $type: "string" }
    }
})


for (let i = 0; i < num_urls; i += bulk_load_count) {
    pages = []
    for (let j = 0; j < bulk_load_count; j++) {
        pages.push({
            _id: NumberInt(i + j),
            "schema": NumberInt(schema_version)
        })
    }
    db.unusedPagesIDs.insertMany(pages)
    pages = []
}

