use chrono::{DateTime, Local};

mod types;
mod fetch;

use types::{Condition, Forecast, TargetDate};

fn main() {
    println!("Hello, world! {:?}", Forecast::new(Local::now(), TargetDate::Now, 20.0, Condition::Cloudy, "Berlin", "accuweather"));
}

// Every hour:
//     For city in all interesting cities:
//         For each provider that works in that city
//             Get all forecasts

// BI:
// Rating of providers among all cities in categories:
//   * Today - closest to the (average on the leadtime 0)
//   * Tomorrow LT 1 (12-36h) - hits the target the best in (average on the leadtime 0)
//   * Tomorrow LT   - hits the target the best in (average on the leadtime 0)
// For each 


// Challenges:
// How to treat distance in time? if it is a date? should it match midday?