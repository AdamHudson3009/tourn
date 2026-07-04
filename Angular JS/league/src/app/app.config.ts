import { ApplicationConfig } from '@angular/core';
import { provideRouter, Routes, withComponentInputBinding } from '@angular/router';

const routes: Routes = [
  { path: '', redirectTo: 'league', pathMatch: 'full' }, // Default route
  { path: 'league', loadComponent: () => import('./pages/league/league.component').then(m => m.LeagueComponent) },
  { path: 'tourn/:league_id/:league_name', loadComponent: () => import('./pages/tourn/tourn.component').then(m => m.TournComponent) },
  { path: 'cdp/:league_id/:league_name/:tourn_id/:tourn_name', loadComponent: () => import('./pages/cdp/cdp.component').then(m => m.CdpComponent) },
  { path: 'trn-table/:league_id/:league_name/:tourn_id/:tourn_name', loadComponent: () => import('./pages/trn-table/trn-table.component').then(m => m.TrnTableComponent) },
  { path: 'tbrks/:league_id/:tourn_id/', loadComponent: () => import('./pages/tbrks/tbrks.component').then(m => m.TbrksComponent) },
  { path: 'cmm/:league_id/:tourn_id/', loadComponent: () => import('./pages/cmm/cmm.component').then(m => m.CmmComponent) },
  { path: 'sch-builder/:league_id/:tourn_id/', loadComponent: () => import('./pages/sch-builder/sch-builder.component').then(m => m.SchBuilderComponent) },
  { path: 'lll-consolidated/:league_id/:tourn_id', loadComponent: () => import('./pages/lll-consolidated/lll-consolidated.component').then(m => m.LllConsolidatedComponent) },
  { path: 'grammar/:league_id/:tourn_id', loadComponent: () => import('./pages/grammar/grammar.component').then(m => m.GrammarComponent) },
  { path: 'trnsch/:league_id/:league_name/:tourn_id/:tourn_name', loadComponent: () => import('./pages/trnsch/trnsch.component').then(m => m.TrnschComponent) },
  { path: 'sch/:league_id/:league_name/:tourn_id/:tourn_name/:rnd/:maxRnd', loadComponent: () => import('./pages/sch/sch.component').then(m => m.SchComponent) },
  { path: 'sch-res/:league_id/:league_name/:tourn_id/:tourn_name/:rnd/:maxRnd/:resStndOoo', loadComponent: () => import('./pages/sch-res/sch-res.component').then(m => m.SchResComponent) },
];

export const appConfig: ApplicationConfig = {
  providers: [provideRouter(routes, withComponentInputBinding())]
};

