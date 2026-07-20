import { Picker } from "@react-native-picker/picker";

interface Props {
  value: string;
  onChange: (value: string) => void;
}

export default function StatusFilter({ value, onChange }: Props) {
  return (
    <Picker selectedValue={value} onValueChange={onChange}>
      <Picker.Item label="All Status" value="" />
      <Picker.Item label="Todo" value="todo" />
      <Picker.Item label="In Progress" value="in_progress" />
      <Picker.Item label="Done" value="done" />
    </Picker>
  );
}
