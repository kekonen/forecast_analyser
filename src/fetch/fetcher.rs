use crate::Forecast;
use std::error;

pub trait ForecastFetcher {
    fn get_forecast(city: &str) -> Result<Forecast, Box<dyn error::Error>>;
    fn get_available_cities() -> Vec<String>;
}