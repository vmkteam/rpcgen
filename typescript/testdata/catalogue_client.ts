/* eslint-disable */
export interface ICampaign {
  id: number | null,
  groups: Array<IGroup> | null
}

export class Campaign implements ICampaign {
  static entityName = "campaign";

  id: number | null = null;
  groups: Array<IGroup> | null = null;
}

export interface ICatalogueFirstParams {
  groups: Array<IGroup> | null
}

export class CatalogueFirstParams implements ICatalogueFirstParams {
  static entityName = "cataloguefirstparams";

  groups: Array<IGroup> | null = null;
}

export interface ICatalogueSecondParams {
  campaigns: Array<ICampaign> | null
}

export class CatalogueSecondParams implements ICatalogueSecondParams {
  static entityName = "cataloguesecondparams";

  campaigns: Array<ICampaign> | null = null;
}

export interface IGroup {
  id: number | null,
  title: string | null,
  nodes: Array<IGroup> | null,
  groups: Array<IGroup> | null,
  child?: IGroup | null,
  sub: ISubGroup | null
}

export class Group implements IGroup {
  static entityName = "group";

  id: number | null = null;
  title: string | null = null;
  nodes: Array<IGroup> | null = null;
  groups: Array<IGroup> | null = null;
  child?: IGroup | null = null;
  sub: ISubGroup | null = null;
}

export interface ISubGroup {
  id: number | null,
  title: string | null,
  nodes: Array<IGroup> | null
}

export class SubGroup implements ISubGroup {
  static entityName = "subgroup";

  id: number | null = null;
  title: string | null = null;
  nodes: Array<IGroup> | null = null;
}

export const factory = (send: any) => ({
  catalogue: {
    first(params: ICatalogueFirstParams): Promise<boolean> {
      return send('catalogue.First', params)
    },
    second(params: ICatalogueSecondParams): Promise<boolean> {
      return send('catalogue.Second', params)
    },
    third(): Promise<object> {
      return send('catalogue.Third')
    }
  }
});
