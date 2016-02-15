package main

type ConfigurationEntry struct {
    StringValue string
    IntValue int
}

var appConfiguration = map[string] ConfigurationEntry {
    //"website_root": "http://www.coshx.com",
    "website_root": ConfigurationEntry{"http://localhost:4567/", 0 },
    //"blog_homepage": "http://www.coshx.com/blog/",
    "blog_homepage": ConfigurationEntry{ "http://localhost:4567/blog/", 0 },
    "min_word_length": ConfigurationEntry { "3", 3 },
    "saved_data_filename": ConfigurationEntry { "saved_data.json", 0 },
}