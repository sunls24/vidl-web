function Badge({ title, value }: { title: string; value: string }) {
  return (
    <div>
      <span className="font-medium">{title}: </span>
      <span className="rounded-xl bg-secondary px-2 py-0.5">{value}</span>
    </div>
  );
}

export default Badge;
