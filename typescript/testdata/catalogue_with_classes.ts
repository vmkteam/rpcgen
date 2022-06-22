/* eslint-disable */
export interface ICampaign {
  id: number,
  groups: Array<IGroup>
}

export interface ICatalogueFirstParams {
  groups: Array<IGroup>
}

export interface ICatalogueSecondParams {
  campaigns: Array<ICampaign>
}

export interface IGroup {
  id: number,
  title: string,
  nodes: Array<IGroup>,
  groups: Array<IGroup>,
  child?: IGroup,
  sub: ISubGroup
}

export interface ISubGroup {
  id: number,
  title: string,
  nodes: Array<IGroup>
}

export class Campaign implements ICampaign {
  static entityName = "campaign";

  id: number = 0;
  groups: Array<IGroup> = null;
}

export class CatalogueFirstParams implements ICatalogueFirstParams {
  static entityName = "cataloguefirstparams";

  groups: Array<IGroup> = null;
}

export class CatalogueSecondParams implements ICatalogueSecondParams {
  static entityName = "cataloguesecondparams";

  campaigns: Array<ICampaign> = null;
}

export class Group implements IGroup {
  static entityName = "group";

  id: number = 0;
  title: string = null;
  nodes: Array<IGroup> = null;
  groups: Array<IGroup> = null;
  child?: IGroup = null;
  sub: ISubGroup = null;
}

export class SubGroup implements ISubGroup {
  static entityName = "subgroup";

  id: number = 0;
  title: string = null;
  nodes: Array<IGroup> = null;
}

export const factory = (send: any) => ({
  catalogue: {
    first(params: ICatalogueFirstParams): Promise<boolean> {
      return send('catalogue.First', params)
    },
    second(params: ICatalogueSecondParams): Promise<boolean> {
      return send('catalogue.Second', params)
    },
    third(): Promise<ICampaign> {
      return send('catalogue.Third')
    }
  }
})
