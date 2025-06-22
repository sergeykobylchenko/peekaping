import React from "react";
import { useTimezone } from "../context/timezone-context";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "./ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { getTimezoneOffsetLabel, sortedTimezones } from "@/lib/timezones";

const TimezoneSelector: React.FC = () => {
  const { timezone, setTimezone } = useTimezone();
  const [open, setOpen] = React.useState(false);
  const [search, setSearch] = React.useState("");

  // Filter by search
  const filteredTimezones = React.useMemo(() => {
    if (!search) return sortedTimezones;
    return sortedTimezones.filter(
      (tz) =>
        tz.toLowerCase().includes(search.toLowerCase()) ||
        getTimezoneOffsetLabel(tz).toLowerCase().includes(search.toLowerCase())
    );
  }, [search]);

  const selectedLabel = timezone
    ? `${timezone} (${getTimezoneOffsetLabel(timezone)})`
    : "Select timezone...";

  return (
    <div className="flex flex-col gap-1">
      <label htmlFor="timezone-combobox" className="text-sm font-medium mb-1">
        Timezone
      </label>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="min-w-[260px] justify-between"
            id="timezone-combobox"
          >
            {selectedLabel}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="min-w-[260px] p-0">
          <Command>
            <CommandInput
              placeholder="Search timezone..."
              value={search}
              onValueChange={setSearch}
              className="h-9"
            />
            <CommandList>
              <CommandEmpty>No timezone found.</CommandEmpty>
              <CommandGroup>
                {filteredTimezones.map((tz) => (
                  <CommandItem
                    key={tz}
                    value={tz}
                    onSelect={() => {
                      setTimezone(tz);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        timezone === tz ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {tz}{" "}
                    <span className="ml-2 text-muted-foreground">
                      ({getTimezoneOffsetLabel(tz)})
                    </span>
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
};

export default TimezoneSelector;
