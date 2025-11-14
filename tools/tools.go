package tools

type GetWeatherParams struct {
	Location string `json:"location"`
}

func GetWeather(params GetWeatherParams) (string, error) {
	if params.Location == "Philadelphia" {
		return "sunny", nil
	}

	return "rainy", nil
}
