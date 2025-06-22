import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useFormContext } from "react-hook-form";
import Timezone from "./timezone";
import StartEndDateTime from "./start-end-date-time";
import cronstrue from "cronstrue";

const CronExpressionForm = () => {
  const form = useFormContext();

  return (
    <>
      <FormField
        control={form.control}
        name="cron"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Cron Expression</FormLabel>
            <FormDescription>
              {cronstrue.toString(field.value, {
                throwExceptionOnParseError: false
              })}
            </FormDescription>
            <FormControl>
              <Input placeholder="30 3 * * *" {...field} />
            </FormControl>
            <FormDescription>Enter a valid cron expression</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="duration"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Duration (Minutes)</FormLabel>
            <FormControl>
              <Input
                type="number"
                min="1"
                step="1"
                {...field}
                onChange={(e) => field.onChange(parseInt(e.target.value) || 0)}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <Timezone />

      <div className="space-y-4">
        <FormLabel>Effective Date Range (Optional)</FormLabel>
        <StartEndDateTime />
      </div>
    </>
  );
};

export default CronExpressionForm;
