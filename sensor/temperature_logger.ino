#include <ESP8266HTTPClient.h>
#include <ESP8266WiFi.h>
#include <DHT.h>
#include <DHT_U.h>

#define DHTPIN D4
#define DHTTYPE DHT11

const char* ssid = "2ABC-2,4";
const char* password = "cajsrumem";

DHT dht(DHTPIN, DHTTYPE);

const char* serverName = "http://192.168.0.200:6969";
unsigned long lastTime = 0;
unsigned long timerDelay = 10000;

void setup() {
	WiFi.begin(ssid, password);
	Serial.begin(115200);
	dht.begin();
	while (WiFi.status() != WL_CONNECTED) {
		delay(1000);
		Serial.println("Connecting to WiFi..");
	}
	Serial.println("Connected to the WiFi network");
	Serial.println(WiFi.localIP());
}

void loop () {
	float h = dht.readHumidity();
	float t = dht.readTemperature();

	if (isnan(h) || isnan(t)) {
		Serial.println("Failed to read from DHT sensor!");
		return;
	}

	if ((millis() - lastTime) > timerDelay) {
		if(WiFi.status() == WL_CONNECTED) {
			WiFiClient client;
			HTTPClient http;

			http.begin(client, serverName);

			http.addHeader("Content-Type", "application/json");
			int httpResponseCode = http.POST("{\"temperature\": " + String(t) + ", \"humidity\": " + String(h) + "}");

			Serial.print("HTTP Response code: ");
			Serial.println(httpResponseCode);

			http.end();
		}
		else {
			Serial.println("WiFi Disconnected");
		}
		lastTime = millis();
	}
}
			

