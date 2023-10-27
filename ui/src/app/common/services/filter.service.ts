import { Injectable } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { FilterOption } from 'app/common/interfaces/orb/filter-option';
import { combineLatest, Observable, ReplaySubject } from 'rxjs';
import { filter, map, shareReplay, tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class FilterService {
  private _filters: FilterOption[];

  private filters: ReplaySubject<FilterOption[]>;

  private activeRoute$: Observable<string>;

  public searchName: string;

  constructor(private router: Router) {
    this.activeRoute$ = this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map((event: NavigationEnd) => event.urlAfterRedirects),
      tap((route) => this.loadFromRoute(route)),
      shareReplay(),
    );
    this.filters = new ReplaySubject();
    this._filters = [];

    this.activeRoute$.subscribe();
  }

  private loadFromRoute(route: string) {
    const storedFilters = window.sessionStorage.getItem(route) || '[]';
    const storedFiltersList = JSON.parse(storedFilters);
    let found = false;

    storedFiltersList.map((filter) => {
      if (filter.name === 'Name' && filter.prop === 'name') {
        this.searchName = filter.param;
        found = true;
      }
    });
    if (!found) {
      this.searchName = '';
    }

    this.resetFilters(JSON.parse(storedFilters));
  }

  private saveToRoute(route: string) {
    const filtersToStore = this._filters.map((_filter) => ({
      name: _filter.name,
      prop: _filter.prop,
      param: _filter.param,
      extra: _filter?.exact,
    }));
    window.sessionStorage.setItem(route, JSON.stringify(filtersToStore));
  }

  getFilters() {
    return this.filters;
  }

  resetFilters(filters: FilterOption[]) {
    this._filters = filters;
    this.filters.next(this._filters);
  }

  cleanFilters() {
    this.commitFilter([]);
  }

  private commitFilter(filters: FilterOption[]) {
    this._filters = filters;
    this.filters.next(this._filters);
    this.saveToRoute(this.router.url);
  }

  addFilter(filterToAdd: FilterOption) {
    this.commitFilter([...this._filters, filterToAdd]);
  }

  removeFilter(index: number) {
    if (index >= 0 && index < this._filters.length) {
      const copy = [...this._filters];
      copy.splice(index, 1);
      this.commitFilter(copy);
    }
  }

  findAndRemove(param: any, name: string) {
    this.removeFilter(this._filters.findIndex((f) => f.param === param && f.name === name && f));
  }

  findFilter(param: any, name: string) {
    return this._filters.findIndex((f) => f.param === param && f.name === name && f);
  }

  // make a decorator out of this?
  createFilteredList() {
    return (
      itemsList: Observable<any[]>,
      filtersList: Observable<FilterOption[]>,
      filterOptions: FilterOption[],
    ) => {
      return combineLatest([itemsList, filtersList]).pipe(
        map(([agents, _filters]) => {
          let filtered = agents;
          _filters.forEach((_filter) => {
            filtered = filtered.filter((value) => {
              const paramValue = _filter.param;
              const filterDef = filterOptions.find(
                (_item) => _item.name === _filter.name,
              );
              if (filterDef) {
                const {filter: filterFn, prop, exact} = filterDef;
                const result =
                !!filterFn && filterFn(value, prop, paramValue, exact);
                return result;
              }
            });
          });

          return filtered;
        }),
      );
    };
  }
}
