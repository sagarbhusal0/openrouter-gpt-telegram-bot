package config

import (
    "github.com/sashabaranov/go-openai"
    "github.com/spf13/viper"
    "log"
    "openrouter-gpt-telegram-bot/lang"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    TelegramBotToken   string
    OpenAIApiKey       string
    Model              ModelParameters
    MaxTokens          int
    BotLanguage        string
    OpenAIBaseURL      string
    SystemPrompt       string
    BudgetPeriod       string
    GuestBudget        float64
    UserBudget         float64
    AdminChatIDs       []int64
    AllowedUserChatIDs []int64
    MaxHistorySize     int
    MaxHistoryTime     int
    Vision            string
    VisionPrompt      string
    VisionDetails     string
    StatsMinRole      string
    Lang              string
}

type ModelParameters struct {
    Type              string
    ModelName         string
    ModelReq          openai.ChatCompletionRequest
    FrequencyPenalty  float64
    MinP              float64
    PresencePenalty   float64
    RepetitionPenalty float64
    Temperature       float64
    TopA              float64
    TopK              float64
    TopP              float64
}

// getStrAsIntList converts a comma-separated string of numbers to []int64
func getStrAsIntList(envKey string) []int64 {
    str := os.Getenv(envKey)
    if str == "" {
        return []int64{}
    }

    strList := strings.Split(str, ",")
    var intList []int64

    for _, s := range strList {
        if s == "" {
            continue
        }
        i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
        if err != nil {
            log.Printf("Warning: could not parse %s as int64: %v", s, err)
            continue
        }
        intList = append(intList, i)
    }
    return intList
}

func Load() (*Config, error) {
    // Set default values
    viper.SetDefault("MAX_TOKENS", 2000)
    viper.SetDefault("TEMPERATURE", 1)
    viper.SetDefault("TOP_P", 0.7)
    viper.SetDefault("BASE_URL", "https://api.openai.com/v1")
    viper.SetDefault("BUDGET_PERIOD", "monthly")
    viper.SetDefault("MAX_HISTORY_SIZE", 10)
    viper.SetDefault("MAX_HISTORY_TIME", 60)
    viper.SetDefault("LANG", "en")

    // Environment variables
    config := &Config{
        TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
        OpenAIApiKey:     os.Getenv("API_KEY"),
        Model: ModelParameters{
            Type:        os.Getenv("TYPE"),
            ModelName:   os.Getenv("MODEL"),
            Temperature: getEnvFloat("TEMPERATURE", 1.0),
            TopP:        getEnvFloat("TOP_P", 0.7),
            FrequencyPenalty: getEnvFloat("FREQUENCY_PENALTY", 0),
            PresencePenalty:  getEnvFloat("PRESENCE_PENALTY", 0),
            MinP:            getEnvFloat("MIN_P", 0),
            RepetitionPenalty: getEnvFloat("REPETITION_PENALTY", 1),
            TopA:            getEnvFloat("TOP_A", 0),
            TopK:            getEnvFloat("TOP_K", 0),
        },
        MaxTokens:          getEnvInt("MAX_TOKENS", 2000),
        OpenAIBaseURL:      getEnvString("BASE_URL", "https://api.openai.com/v1"),
        SystemPrompt:       os.Getenv("ASSISTANT_PROMPT"),
        BudgetPeriod:       getEnvString("BUDGET_PERIOD", "monthly"),
        GuestBudget:        getEnvFloat("GUEST_BUDGET", 0),
        UserBudget:         getEnvFloat("USER_BUDGET", 0),
        AdminChatIDs:       getStrAsIntList("ADMIN_IDS"),
        AllowedUserChatIDs: getStrAsIntList("ALLOWED_USER_IDS"),
        MaxHistorySize:     getEnvInt("MAX_HISTORY_SIZE", 10),
        MaxHistoryTime:     getEnvInt("MAX_HISTORY_TIME", 60),
        Vision:             os.Getenv("VISION"),
        VisionPrompt:       os.Getenv("VISION_PROMPT"),
        VisionDetails:      os.Getenv("VISION_DETAIL"),
        StatsMinRole:       getEnvString("STATS_MIN_ROLE", "user"),
        Lang:               getEnvString("LANG", "en"),
    }

    // Validate required configurations
    if config.TelegramBotToken == "" {
        return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
    }
    if config.OpenAIApiKey == "" {
        return nil, fmt.Errorf("API_KEY is required")
    }
    if config.BudgetPeriod == "" {
        return nil, fmt.Errorf("BUDGET_PERIOD is required")
    }

    // Load language
    language := lang.Translate("language", config.Lang)
    if language == "" {
        log.Printf("Warning: Language '%s' not found, defaulting to 'en'", config.Lang)
        config.Lang = "en"
    }

    return config, nil
}

// Helper functions to get environment variables with default values
func getEnvString(key string, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func getEnvInt(key string, defaultValue int) int {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    intValue, err := strconv.Atoi(value)
    if err != nil {
        log.Printf("Warning: could not parse %s as int: %v, using default value %d", key, err, defaultValue)
        return defaultValue
    }
    return intValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    floatValue, err := strconv.ParseFloat(value, 64)
    if err != nil {
        log.Printf("Warning: could not parse %s as float64: %v, using default value %f", key, err, defaultValue)
        return defaultValue
    }
    return floatValue
}