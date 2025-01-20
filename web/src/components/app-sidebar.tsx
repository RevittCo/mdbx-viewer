import { CopyMinus, Database } from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import { Tree, Folder, File, CollapseButton } from "@/components/tree-view-api";
import { DataSource, TreeElement } from "@/types";
import { useApi } from "@/api/api";
import { useEffect, useState } from "react";

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  handleOpenDbSource: (source: string, name: string) => Promise<void>;
};

export function AppSidebar({ handleOpenDbSource, ...props }: AppSidebarProps) {
  const api = useApi();
  const [dataSource, setDataSource] = useState<DataSource | undefined>(undefined);

  const fetchDataSource = async () => {
    const res = await api.buckets.getDataSource();
    if (res.data && res.status === 200) {
      setDataSource(res.data);
    }
  };

  useEffect(() => {
    fetchDataSource();
  }, []);

  const fileClick = (path: string, name: string) => {
    handleOpenDbSource(path, name);
  };

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="#">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Database className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-semibold">MDBX Viewer</span>
                  <span className="">v1.0.0</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            <Tree className="rounded-md bg-background p-2" initialSelectedId="21" elements={dataSource?.treeElements}>
              {dataSource && (
                <>
                  <div className="w-full flex items-start justify-between gap-1">
                    <div className="text-sm" style={{ wordBreak: "break-all" }}>
                      {dataSource?.source || "Data Directory"}
                    </div>
                    <CollapseButton
                      elements={dataSource?.treeElements}
                      className="p-0 h-auto w-auto m-0 relative rounded-none hover:bg-transparent"
                    >
                      <CopyMinus className="transform -scale-x-100" />
                    </CollapseButton>
                  </div>
                  <TreeElementRenderer elements={dataSource.treeElements} fileClick={fileClick} />
                </>
              )}
            </Tree>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  );
}

const TreeElementRenderer: React.FC<{ elements: TreeElement[], fileClick: (path: string, name: string) => void }> = ({ elements, fileClick }) => {
  return (
    <>
      {elements.map((element) => {
        if (element.children && element.children.length > 0) {
          return (
            <Folder key={element.fullPath} element={element.name} value={element.fullPath}>
              <TreeElementRenderer elements={element.children} fileClick={fileClick} />
            </Folder>
          );
        }
        return (
          <File key={element.fullPath} value={element.fullPath}>
            <div onClick={() => fileClick(element.fullPath, element.name)}>{element.name}</div>
          </File>
        );
      })}
    </>
  );
};
