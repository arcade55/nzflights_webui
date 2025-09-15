package components

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
)

func FlightCardComponent(flightValue nzflights.FlightValue) htma.Element {
	f := flightValue.Flight

	return htma.FlightCard().
		FlightNumberAttr(f.IdentIATA).
		AirlineNameAttr(getAirlineName(f.Operator)).
		OriginIataAttr(f.OriginIATA).
		OriginCityAttr(f.OriginCity).
		DestIataAttr(f.DestinationIATA).
		DestCityAttr(f.DestinationCity).
		GateAttr(f.GateOrigin).
		DepartureTimeAttr(f.ScheduledOut).
		ArrivalTimeAttr(f.ScheduledIn).
		StatusTextAttr(f.Status)

}

/*
	func formatTime(isoString string) string {
		t, err := time.Parse(time.RFC3339, isoString)
		if err != nil {
			return "??:??"
		}
		return t.Format("03:04 PM")
	}
*/
func getAirlineName(operator string) string {
	// In a real app, this would be more robust
	if operator == "ANZ" {
		return "Air New Zealand"
	}
	if operator == "JST" {

		return "Jetstar"
	}

	return ""
}
