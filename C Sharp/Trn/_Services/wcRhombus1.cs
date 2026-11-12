using Dapper;
using Trn.Data;
using Trn.Dtos;
using Microsoft.Extensions.Configuration;

namespace Trn.Services;

public static class WCRhombus1
{
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());

    public static (string, List<int>, List<int>) ProcessWCRhombus1(int trnId, List<int> plyrsIds, int wcBase, List<int> wcWns, int currRnd, int rowRnd) {
        MatchResponse mtchResponse = wc_rhombus1(plyrsIds, wcWns, wcBase);
        
        if (rowRnd < currRnd) {
            if (!WCCommon.wc_insert(trnId, currRnd, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids)) {
                return ("Rhombus1: Failed to insert matches.", new List<int>(), new List<int>());
            }
            return ("Rhombus1: Matches Inserted", new List<int>(), new List<int>());
        }
        bool isValid;

        (isValid, wcWns) = WCCommon.wc_validate(trnId, currRnd, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids);
        if (!isValid) {
            return ("Rhombus1: Matches not consistent with flow.", new List<int>(), new List<int>());
        }
        for (int indx = 0; indx < wcWns.Count; indx++) {
            if (wcWns[indx] < 0)
                return ("Rhombus1: Matches not completed.", new List<int>(), new List<int>());
        }

        List<int> remainPlyrs = new List<int>();
        remainPlyrs.AddRange(mtchResponse.Plyrs1Ids);
        remainPlyrs.AddRange(mtchResponse.Plyrs2Ids);
        return ("", wcWns, remainPlyrs);
    }

    private static MatchResponse wc_rhombus1(List<int> plyrsIds, List<int> wcWns, int wcBase) {
        var plyrTmp = new List<int>();
 
        string wcWnsStr = $",{string.Join(",", wcWns)},";
        for (int indx = 0; indx < plyrsIds.Count; indx++) {
            if (indx < (wcBase - wcWns.Count) || (wcWnsStr.Contains($",{plyrsIds[indx]},"))) {
                plyrTmp.Add(plyrsIds[indx]);
            }
        }
        return WCCommon.wc_round_setup(plyrTmp);
    }
}