import { View, Button, Text, StyleSheet } from "react-native";

interface Props {
  page: number;
  total: number;
  limit: number;
  onPrevious: () => void;
  onNext: () => void;
}

export default function Pagination({
  page,
  total,
  limit,
  onPrevious,
  onNext,
}: Props) {
  const totalPages = Math.ceil(total / limit);

  return (
    <View style={styles.container}>
      <Button title="Previous" onPress={onPrevious} disabled={page === 1} />

      <Text>
        {page} / {totalPages || 1}
      </Text>

      <Button title="Next" onPress={onNext} disabled={page >= totalPages} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    margin: 10,
  },
});
