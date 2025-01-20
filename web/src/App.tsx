import { useApi } from "./api/api";
import { AppSidebar } from "@/components/app-sidebar";
import { Separator } from "@/components/ui/separator";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useEffect, useState } from "react";
import { KeyValuePair } from "./types";
import DataTable from "./components/data-table";
import { File, RefreshCcw } from "lucide-react";
import { Input } from "./components/ui/input";
import { Button } from "./components/ui/button";

function App() {
  const api = useApi();
  const [source, setSource] = useState<string | undefined>(undefined);
  const [sources, setSources] = useState<string[]>([]);
  const [bucketValue, setBucketValue] = useState<string | undefined>(undefined);
  const [pages, setPages] = useState<KeyValuePair[]>([]);

  const [limit, setLimit] = useState<string>("100");
  const [offset, setOffset] = useState<string>("0");

  const handleOpenDbSource = async (source: string, name: string) => {
    const res = await api.buckets.openDBSource(source);
    if (res.data && res.status === 200) {
      setSource(name);
      setSources(res.data);
    }
  };

  const handleReadPage = async (name: string, num: string, len: string) => {
    const res = await api.buckets.readBucket(name, num, len);
    if (res.data && res.status === 200) {
      setPages(res.data);
    } else {
      setPages([]);
    }
  };

  useEffect(() => {
    if (bucketValue) {
      handleReadPage(bucketValue, offset, limit);
    }
  }, [bucketValue]);

  return (
    <SidebarProvider>
      <AppSidebar handleOpenDbSource={handleOpenDbSource} />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b w-full">
          <div className="flex items-center gap-2 px-3 w-full">
            <SidebarTrigger />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <div className="flex items-center justify-between w-full gap-2">
              {source ? (
                <>
                  <div className="whitespace-nowrap h-9 flex items-center justify-center pl-2 pr-3 text-muted-foreground">
                    <File className="h-4" /> {source}
                  </div>
                  {sources.length > 0 && (
                    <Select value={bucketValue} onValueChange={setBucketValue}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select bucket name" />
                      </SelectTrigger>
                      <SelectContent>
                        {sources.map((source, index) => (
                          <SelectItem key={index} value={source}>
                            {source}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                  {bucketValue && (<Button className="px-3 text-muted-foreground" variant={"outline"} onClick={() => handleReadPage(bucketValue, offset, limit)}><RefreshCcw /></Button>)}
                  {bucketValue && (
                    <div className="flex gap-2 items-center">
                      <div className="flex items-center gap-2 text-muted-foreground">
                        Limit
                        <Input type="number" value={limit} onChange={(e) => setLimit(e.target.value)} className="w-20" />
                      </div>
                      <div className="flex items-center gap-2 text-muted-foreground">
                        Offset
                        <Input type="number" value={offset} onChange={(e) => setOffset(e.target.value)} className="w-20" />
                      </div>
                    </div>
                  )}
                  {bucketValue && (
                    <div className="whitespace-nowrap text-muted-foreground px-2">
                      {pages.length} Result{pages.length === 0 || pages.length > 1 ? "'s" : ""}
                    </div>
                  )}
                </>
              ) : (
                <div className="whitespace-nowrap text-muted-foreground">No source loaded</div>
              )}
            </div>
          </div>
        </header>
        <div className="p-4 w-full flex overflow-auto h-full">
          <div className="flex-grow">
            <DataTable pages={pages} />
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

export default App;
