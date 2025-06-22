import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import StartEndTime from "./start-end-time";
import Timezone from "./timezone";
import StartEndDateTime from "./start-end-date-time";
import { Checkbox } from "@/components/ui/checkbox";
import { useFormContext } from "react-hook-form";

const WEEKDAYS = [
  { id: "0", label: "Sun", value: 0 },
  { id: "1", label: "Mon", value: 1 },
  { id: "2", label: "Tue", value: 2 },
  { id: "3", label: "Wed", value: 3 },
  { id: "4", label: "Thu", value: 4 },
  { id: "5", label: "Fri", value: 5 },
  { id: "6", label: "Sat", value: 6 },
];

const RecurringWeekdayForm = () => {
  const form = useFormContext();

  return (
    <>
      <FormField
        control={form.control}
        name="weekdays"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Day of Week</FormLabel>
            <div className="flex gap-4">
              {WEEKDAYS.map((weekday) => (
                <FormItem
                  key={weekday.id}
                  className="flex flex-col items-center space-y-0.5"
                >
                  <FormLabel className="text-xs text-gray-600">{weekday.label}</FormLabel>
                  <FormControl>
                    <Checkbox
                      checked={field.value?.includes(weekday.value)}
                      onCheckedChange={(checked) => {
                        const current = field.value || [];
                        if (checked) {
                          field.onChange([...current, weekday.value]);
                        } else {
                          field.onChange(
                            current.filter((v: number) => v !== weekday.value)
                          );
                        }
                      }}
                    />
                  </FormControl>
                </FormItem>
              ))}
            </div>
            <FormMessage />
          </FormItem>
        )}
      />

      <StartEndTime />
      <Timezone />

      <div className="space-y-4">
        <FormLabel>Effective Date Range (Optional)</FormLabel>
        <StartEndDateTime />
      </div>
    </>
  );
};

export default RecurringWeekdayForm;
