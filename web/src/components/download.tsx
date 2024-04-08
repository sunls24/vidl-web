import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog.tsx";
import { RotateCw } from "lucide-react";
import { Progress } from "@/components/ui/progress.tsx";
import { ReactNode } from "react";

function Download({
  open,
  onChange,
  progress,
  trigger,
}: {
  open: boolean;
  onChange: (open: boolean) => void;
  progress: number;
  trigger: ReactNode;
}) {
  return (
    <AlertDialog open={open} onOpenChange={onChange}>
      <AlertDialogTrigger asChild>{trigger}</AlertDialogTrigger>
      <AlertDialogContent>
        <div className="flex items-center justify-center gap-2">
          <RotateCw className="animate-spin" />
          <span className="font-medium">
            The video file is being
            <span className="mx-1 font-bold underline underline-offset-4">
              {progress === 0 ? "prepared" : "downloaded"}
            </span>
            ⌛
          </span>
        </div>
        {progress !== 0 && <Progress value={progress} />}
      </AlertDialogContent>
    </AlertDialog>
  );
}

export default Download;
