import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import StartEndTime from "./start-end-time";
import Timezone from "./timezone";
import StartEndDateTime from "./start-end-date-time";
import { Input } from "@/components/ui/input";
import { useFormContext } from "react-hook-form";

const RecurringIntervalForm = () => {
  const form = useFormContext();

  return (
    <>
      <FormField
        control={form.control}
        name="intervalDay"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Interval (Run once every day)</FormLabel>
            <FormControl>
              <Input
                type="number"
                min="1"
                max="3650"
                step="1"
                {...field}
                onChange={(e) => field.onChange(parseInt(e.target.value) || 1)}
              />
            </FormControl>
            <FormDescription>
              {field.value &&
                field.value >= 1 &&
                `Every ${field.value} day${field.value > 1 ? "s" : ""}`}
            </FormDescription>
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

export default RecurringIntervalForm;
