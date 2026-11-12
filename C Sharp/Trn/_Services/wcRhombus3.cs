using Dapper;
using Trn.Data;
using Trn.Dtos;
using Microsoft.Extensions.Configuration;

namespace Trn.Services;

public static class WCRhombus3
{
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());

    public static string ProcessWCRhombus3(int trnId, List<int> plyrsIds, List<int> wcWns, int currRnd, int rowRnd) {
        Dictionary<int, int> plyrsRecDict = new Dictionary<int, int>();
        plyrsRecDict = GetPlyrsRecDict(trnId, plyrsIds, currRnd);

        Console.WriteLine($"Player Records: {string.Join(",", plyrsRecDict.Select(kvp => kvp.Key + ":" + kvp.Value))}");

        MatchResponse mtchResponse = wc_rhombus3(plyrsIds, wcWns, plyrsRecDict);
        
        if (rowRnd < currRnd) {
            if (!WCCommon.wc_insert(trnId, currRnd, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids, mtchResponse.Plyrs1Ids.Count/2)) {
                return ("Rhombus3: Failed to insert matches.");
            }
            return ("Rhombus3: Matches Inserted");
        }
        bool isValid;

        (isValid, wcWns) = WCCommon.wc_validate(trnId, currRnd, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids);
        if (!isValid) {
            return ("Rhombus3: Matches not consistent with flow.");
        }
        for (int indx = 0; indx < wcWns.Count; indx++) {
            if (wcWns[indx] < 0)
                return ("Rhombus3: Matches not completed.");
        }

        return ("");
    }

    private static Dictionary<int, int> GetPlyrsRecDict(int trnId, List<int> plyrsIds, int currRnd) {
        using var conn = _db.Create();
        var rows = conn.Query(
            @"SELECT plyrtm_id, max(rec) AS rec
            FROM gms_plyrtm
            WHERE gms_id IN (SELECT id FROM gms WHERE tourn_id = @TrnId 
            and coalesce(has_stnd,'Y') = 'N' and rnd < @CurrRnd)
            AND plyrtm_id IN @PlyrsIds
            GROUP BY plyrtm_id",
            new { TrnId = trnId, PlyrsIds = plyrsIds, CurrRnd = currRnd }).ToList();

        var plyrsRecDict = new Dictionary<int, int>();
        foreach (var row in rows) {
            plyrsRecDict[row.plyrtm_id] = row.rec;
        }
        return plyrsRecDict;
    }

    private static MatchResponse wc_rhombus3(List<int> plyrsIds, List<int> wcWns, Dictionary<int, int> plyrsRecDict) {
        var plyrsNext = new List<int>();
        string wcWnsStr = $",{string.Join(",", wcWns)},";

        for (int indx = 0; indx < plyrsIds.Count; indx++) {
            if (wcWnsStr.Contains($",{plyrsIds[indx]},")) {
                if (plyrsRecDict[plyrsIds[indx]] == 0) {
                    plyrsNext.Add(plyrsIds[indx]);
                } 
            } else if (plyrsRecDict[plyrsIds[indx]] == 1) {
                plyrsNext.Add(plyrsIds[indx]);
            }
        }

        MatchResponse mtchResponseNext = WCCommon.wc_round_setup(plyrsNext);
        return mtchResponseNext;
    }
}