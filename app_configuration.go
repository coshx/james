package main

type ConfigurationEntry struct {
    StringValue string
    IntValue int
    FloatValue float64
}

var appConfiguration = map[string] ConfigurationEntry {
    //"website_root": "http://www.coshx.com",
    "website_root": ConfigurationEntry{"http://localhost:4567", 0, 0 },
    //"blog_homepage": "http://www.coshx.com/blog/",
    "blog_homepage": ConfigurationEntry{ "http://localhost:4567/blog/", 0, 0 },
    "min_word_length": ConfigurationEntry { "3", 3, 0 },
    "saved_data_filename": ConfigurationEntry { "saved_data.json", 0, 0 },
    "compare_words_ratio": ConfigurationEntry { "0.60", 0, 0.60 },
    "headline_coeff": ConfigurationEntry { "15", 15, 0 },
    "author_coeff": ConfigurationEntry { "10", 10, 0 },
    "minimum_weight": ConfigurationEntry { "5.0", 0, 5.0},
}