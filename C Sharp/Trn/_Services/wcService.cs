using Dapper;
using Trn.Data;
using Trn.Dtos;
using Microsoft.Extensions.Configuration;

namespace Trn.Services;

public class WCService
{
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());

    public async Task<OutResponse> ProcessWCAsync(ItemRequest req) {
        if (req.ConfId == null || req.PlyrsConIds == null || 
        req.ConfId.Count != req.PlyrsConIds.Count) {
            return new OutResponse($"ConfId {req.ConfId?.Count} and PlyrsConIds {req.PlyrsConIds?.Count} must be non-null and have the same count.");
        }
        int wcBase = WCBase.GetWCBase(req.TrnId);

        string msg = "";
        bool wcPrelim = false;
        for (int indx = 0; indx < req.ConfId.Count; indx++) {
            if (req.PlyrsConIds[indx].Count > wcBase) {
                wcPrelim = true;
                break;
            }
        }

        for (int indx = 0; indx < req.ConfId.Count; indx++) {
            string tempMsg = ProcessConfWC(req.TrnId, wcBase, req.ConfId[indx], req.PlyrsConIds[indx], wcPrelim);
            if (tempMsg != "") msg += $"req.ConfId[{indx}]: {tempMsg} ";
        }
        return new OutResponse(msg);
    }

    private string ProcessConfWC(int trnId, int wcBase, int confId, List<int> plyrsIds, bool wcPrelim) {
        using var conn = _db.Create();

        var rows = conn.Query(
            @"SELECT MAX(rnd) AS rnd,
             SUM(CASE WHEN COALESCE(plyrtm_wn, -2) < 0 THEN 1 ELSE 0 END) AS remaining
            FROM gms
            WHERE tourn_id = @TrnId AND COALESCE(has_stnd, 'Y') = 'Y'",
            new { TrnId = trnId }).ToList();

        if (rows[0].remaining > 0) return "Matches assigned not all played.";

        int maxRnd = rows[0].rnd;
        rows = conn.Query(
            @"SELECT COALESCE(MAX(g.rnd), 0) AS rnd
            FROM gms g
            inner join gms_plyrtm gp on gp.gms_id = g.id
            inner join plyrtm p on p.id = gp.plyrtm_id
            inner join divi d on d.id = p.divi_id
            inner join conf c on c.id = d.conf_id
            WHERE g.tourn_id = @TrnId AND COALESCE(has_stnd, 'Y') <> 'Y'
            and c.id = @ConfId",
            new { TrnId = trnId, ConfId = confId }).ToList();

        List<int> wcWns = new List<int>();

        string msg = "";
        int rndRow = 0;
        if (rows.Count > 0) rndRow = (int)rows[0].rnd;

        if (wcPrelim) {
            (msg, wcWns) = WCPrelim.ProcessWCPrelim(trnId, plyrsIds, wcBase, maxRnd, rndRow);
            if (msg != "") return msg;
        }

        int currRnd = maxRnd+1;
        if (wcPrelim) {
            currRnd++;
        }

        (msg, wcWns, plyrsIds) = WCRhombus1.ProcessWCRhombus1(trnId, plyrsIds, wcBase, wcWns, currRnd, rndRow);
        if (msg != "") return msg;
        currRnd++;
        
        (msg, wcWns, plyrsIds) = WCRhombus2.ProcessWCRhombus2(trnId, plyrsIds, wcWns, currRnd, rndRow);
        if (msg != "") return msg;
        currRnd++;
        
        msg = WCRhombus3.ProcessWCRhombus3(trnId, plyrsIds, wcWns, currRnd, rndRow);
        if (msg != "") return msg;
        
        return "WC Done";
    }
}