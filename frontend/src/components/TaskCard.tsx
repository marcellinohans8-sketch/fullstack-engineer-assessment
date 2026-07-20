import { View, Text, Button, StyleSheet } from "react-native";

import { Task } from "../types/task";

type Props = {
  task: Task;
  onEdit: () => void;
};

export default function TaskCard({ task, onEdit }: Props) {
  return (
    <View style={styles.card}>
      <Text style={styles.title}>{task.title}</Text>

      <Text>{task.description}</Text>

      <Text>Status: {task.status}</Text>

      <Text>Assignee: {task.assignee}</Text>

      <View style={styles.buttonContainer}>
        <Button title="Edit" onPress={onEdit} />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: "#fff",
    marginHorizontal: 10,
    marginVertical: 6,
    padding: 15,
    borderRadius: 8,
    elevation: 2,
  },

  title: {
    fontSize: 18,
    fontWeight: "bold",
    marginBottom: 5,
  },

  buttonContainer: {
    marginTop: 10,
  },
});
