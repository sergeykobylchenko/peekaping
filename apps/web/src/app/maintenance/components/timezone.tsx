import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { getTimezoneOffsetLabel, sortedTimezones } from "@/lib/timezones";
import { useFormContext } from "react-hook-form";

const timezoneOptions = [
  { value: "SAME_AS_SERVER", label: "Same as Server Timezone" },
  { value: "UTC", label: "UTC" },
  ...sortedTimezones.map((el) => ({
    value: el,
    label: `${el} (${getTimezoneOffsetLabel(el)})`,
  })),
];

const Timezone = () => {
  const form = useFormContext();

  return (
    <FormField
      control={form.control}
      name="timezone"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Timezone</FormLabel>
          <Select onValueChange={field.onChange} value={field.value}>
            <FormControl>
              <SelectTrigger>
                <SelectValue placeholder="Select timezone" />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              {timezoneOptions.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};

export default Timezone;
