import { useState } from "react";
import { Button } from "@/components/ui/button.tsx";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="mx-auto mt-20 w-fit">
      <h1 className="text-2xl font-medium">Vite + React</h1>
      <div>
        <Button className="my-4" onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </Button>
        <p>
          Edit
          <code className="mx-1 rounded bg-input/80 px-1 py-0.5">
            src/App.tsx
          </code>
          and save to test HMR
        </p>
      </div>
    </div>
  );
}

export default App;
