
const resource: Record<string, string> = {};

resource["CMM"] = "common";
resource["Conf"] = "Conference";
resource["Divi"] = "Division";
resource["trnsch"] = "View Schedule";
resource["Msc"] = "Misc";
resource["Pa"] = "Points Against";
resource["Pf"] = "Points For";
resource["POTW"] = "Player of the Round";
resource["Pt"] = "Net Points";
resource["Rec"] = "Wins (Players)";
resource["sch"] = "Round";
resource["schres"] = "Standings";

resource["Tourns"] = "Tournaments";

function get_resource(key: string, use: boolean = true) {
    if (use && resource[key]) {
        return resource[key];
    }
    return key
}