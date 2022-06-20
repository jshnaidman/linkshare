db_name = DB_NAME
schema_version = SCHEMA_VERSION
num_urls = URL_CAPACITY
bulk_load_count = RUN_ONCE_BULK_LOAD_PAGE_AMOUNT

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
                    bsonType: "string",
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
                    bsonType: "string",
                    description: "Username"
                },
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
                        bsonType: "int",
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
