import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { KeyValuePair } from "@/types";
import { Copy } from "lucide-react";

interface DataTableProps {
  pages: KeyValuePair[];
}

const DataTable = ({ pages }: DataTableProps) => {
  return (
    <div className="rounded-lg border shadow-sm">
      <Table className="w-full">
        <TableHeader className="bg-muted/50">
          <TableRow>
            <TableHead className="py-1 px-1 h-auto text-left text-sm font-semibold w-1/24 border-r"></TableHead>
            <TableHead className="py-1 px-3 h-auto text-left text-sm font-semibold w-11/24 border-r overflow-hidden">
              <p>Key</p>
            </TableHead>
            <TableHead className="py-1 px-3 h-auto text-left text-sm font-semibold w-12/24 overflow-hidden">
              <p>Value</p>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody className="h-full w-full">
          {pages.length > 0 ? (
            pages.map((page, index) => (
              <TableRow key={index} className="w-full">
                <TableCell className="px-3 py-0 text-sm border-r">{index + 1}</TableCell>
                <TableCell className="px-3 text-sm border-r py-0 w-full">
                  <div className="flex flex-col py-2 w-full overflow-auto">
                    <div className="overflow-hidden">{page.key}</div>
                    <div className="flex items-center gap-2">
                      <div className="border h-6 px-2 flex items-center justify-center rounded-md text-muted-foreground whitespace-nowrap">
                        Length {page.key.length}
                      </div>
                      <button
                        className="border h-6 w-6 flex items-center justify-center rounded-md text-muted-foreground hover:bg-muted-foreground/10"
                        onClick={() => navigator.clipboard.writeText(page.key)}
                      >
                        <Copy className="h-3" />
                      </button>
                    </div>
                  </div>
                </TableCell>
                <TableCell className="px-3 py-0 text-sm">
                  <div className="flex flex-col py-2">
                    <div className="overflow-hidden">{page.value}</div>
                    <div className="flex items-center gap-2">
                      <div className="border h-6 px-2 flex items-center justify-center rounded-md text-muted-foreground whitespace-nowrap">
                        Length {page.value.length}
                      </div>
                      <button
                        className="border h-6 w-6 flex items-center justify-center rounded-md text-muted-foreground hover:bg-muted-foreground/10"
                        onClick={() => navigator.clipboard.writeText(page.value)}
                      >
                        <Copy className="h-3" />
                      </button>
                    </div>
                  </div>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={3} className="text-center py-4 text-gray-500">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};

export default DataTable;
