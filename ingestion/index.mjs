import assert from "assert";
import fs from "fs/promises";

const { API_HOST, API_TOKEN } = process.env;
assert(API_HOST, "API_HOST is not set");
assert(API_TOKEN, "API_TOKEN is not set");
assert(URL.canParse(API_HOST), "API_HOST is not a valid URL");

const rawWeatherRecords = await fs.readFile("data/weather.dat", "utf-8");

const weatherRecords = rawWeatherRecords
  .trim()
  .split(/\r?\n/)
  .map((row) => {
    const [date, humidity, temperature] = row.split("\t");
    return {
      date,
      humidity: Number(humidity),
      temperature: Number(temperature),
    };
  });

console.log("🌞 Preparations complete");


// Assignment Note: Could use Promise.all to process records in parallel
console.log(
  `Ready to ingest ${weatherRecords.length} weather records ⛅️🌂...`
);
console.time("Ingestion");
for (const record of weatherRecords) {
  console.log("Processing record:", record);
  try {
    const response = await fetch(`${API_HOST}/weather`, {
      method: "POST",
      headers: { "Content-Type": "application/json", "X-Api-Token": API_TOKEN },
      body: JSON.stringify(record),
    });

    if (!response.ok) {
      throw new Error(`Network error! status: ${response.status}`);
    }
    console.log("Record sent successfully:", record);
  } catch (error) {
    console.error("🌧️ Error sending record:", error);
    console.error("Stop further processing");
    process.exit(1);
  }
}
console.log("🌞 Ingestion complete");
console.timeEnd("Ingestion");
