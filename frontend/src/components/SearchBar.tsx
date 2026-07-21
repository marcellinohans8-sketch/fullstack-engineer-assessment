import { TextInput, StyleSheet } from "react-native";

interface Props {
  value: string;
  onChangeText: (text: string) => void;
}

export default function SearchBar({ value, onChangeText }: Props) {
  return (
    <TextInput
      testID="task-search-input"
      placeholder="Search task..."
      value={value}
      onChangeText={onChangeText}
      style={styles.input}
    />
  );
}

const styles = StyleSheet.create({
  input: {
    backgroundColor: "#fff",
    margin: 10,
    paddingHorizontal: 15,
    paddingVertical: 10,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "#ddd",
  },
});
