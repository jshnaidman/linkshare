db_name = DB_NAME
schema_version = SCHEMA_VERSION
num_urls = 1.2e6
bulk_load_count = num_urls / 10

db.createCollection('free_pages', {
    validator: {
        $jsonSchema: {
            required: ["schema_version"],
            properties: {
                page_id: {
                    bsonType: "int",
                    description: "The id of the pages which aren't being used"
                },
                schema_version: {
                    bsonType: "int"
                }
            }
        }
    }
})

db.createCollection('pages', {
    validator: {
        $jsonSchema: {
            required: ["schema_version"],
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
                    allOf: [{
                        bsonType: "string",
                        maxLength: 2048,
                        description: "A URL. The max length of a URL is 2048 characters"
                    }],
                    maxItems: 200
                },
                user_id: {
                    bsonType: "string",
                    description: "The user_id of whoever owns the page"
                },
                schema_version: {
                    bsonType: "int"
                }
            }
        }
    }
})
db.createCollection('users', {
    validator: {
        $jsonSchema: {
            required: ["username", "schema_version"],
            properties: {
                _id: {
                    bsonType: "string",
                    description: "User ID which comes from an auth service (e.g google)"
                },
                username: {
                    bsonType: "string",
                    description: "chosen username."
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
                    allOf: [{
                        bsonType: "int",
                        description: "The pages owned by the user."
                    }]
                },
                schema_version: {
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
            "schema_version": NumberInt(schema_version)
        })
    }
    db.free_pages.insertMany(pages)
    pages = []
}

// sort them randomly
db.free_pages.aggregate([{ "$sample": { "size": num_urls } }, { "$out": { db: db_name, coll: "free_pages" } }], { allowDiskUse: true })
