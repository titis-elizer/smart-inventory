import { useSelector } from "react-redux";
export default function ItemSelect({ value, onChange }) {
  const items = useSelector((s) => s.inventory.items);

  return (
    <select
      className="input"
      value={value}
      onChange={(e) => onChange(e.target.value)}
    >
      <option value="">Select product</option>

      {items.map((item) => (
        <option key={item.id} value={item.id}>
          {item.name}
        </option>
      ))}
    </select>
  );
}