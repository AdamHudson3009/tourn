using Dapper;
using Trn.Data;
using Trn.Dtos;
using Microsoft.Extensions.Configuration;

namespace Trn.Services;

public static class WCRhombus2
{
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());

    public static (string, List<int>, List<int>) ProcessWCRhombus2(int trnId, List<int> plyrsIds, List<int> wcWns, int currRnd, int rowRnd) {
        MatchResponse mtchResponse = wc_rhombus2(plyrsIds, wcWns);
        
        if (rowRnd < currRnd) {
            if (!WCCommon.wc_insert(trnId, currRnd, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids, mtchResponse.Plyrs1Ids.Count/2)) {
                return ("Rhombus2: Failed to insert matches.", new List<int>(), new List<int>());
            }
            return ("Rhombus2: Matches Inserted", new List<int>(), new List<int>());
        }
        bool isValid;

        (isValid, wcWns) = WCCommon.wc_validate(trnId, currRnd, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids);
        if (!isValid) {
            return ("Rhombus2: Matches not consistent with flow.", new List<int>(), new List<int>());
        }
        for (int indx = 0; indx < wcWns.Count; indx++) {
            if (wcWns[indx] < 0)
                return ("Rhombus2: Matches not completed.", new List<int>(), new List<int>());
        }

        List<int> remainPlyrs = new List<int>();
        remainPlyrs.AddRange(mtchResponse.Plyrs1Ids);
        remainPlyrs.AddRange(mtchResponse.Plyrs2Ids);
        return ("", wcWns, remainPlyrs);
    }

    private static MatchResponse wc_rhombus2(List<int> plyrsIds, List<int> wcWns) {
        var plyrsWn = new List<int>();
        var plyrsLs = new List<int>();
        string wcWnsStr = $",{string.Join(",", wcWns)},";

        for (int indx = 0; indx < plyrsIds.Count; indx++) {
            if (wcWnsStr.Contains($",{plyrsIds[indx]},")) {
                plyrsWn.Add(plyrsIds[indx]);
            } else {
                plyrsLs.Add(plyrsIds[indx]);
            }
        }

        MatchResponse mtchResponseLs = WCCommon.wc_round_setup(plyrsLs);
        MatchResponse mtchResponseWn = WCCommon.wc_round_setup(plyrsWn);
        List<int> resultPlyrs1Ids = mtchResponseLs.Plyrs1Ids.Concat(mtchResponseWn.Plyrs1Ids).ToList();
        List<int> resultPlyrs2Ids = mtchResponseLs.Plyrs2Ids.Concat(mtchResponseWn.Plyrs2Ids).ToList();
        return new MatchResponse(resultPlyrs1Ids, resultPlyrs2Ids);
    }
}