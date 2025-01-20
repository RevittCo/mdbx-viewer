export type KeyValuePair = {
  key: string;
  value: string;
};

export type DataSource = {
  source: string;
  treeElements: TreeElement[];
};

export type TreeElement = {
  id: string;
  name: string;
  fullPath: string;
  isSelectable?: boolean;
  children?: TreeElement[] | null;
};
