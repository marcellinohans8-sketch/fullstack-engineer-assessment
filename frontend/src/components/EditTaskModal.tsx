import { useEffect, useState } from "react";
import { Modal, View, Text, TextInput, Button, StyleSheet } from "react-native";

import { Picker } from "@react-native-picker/picker";

import { Task, TaskStatus } from "../types/task";
import { updateTask } from "../api/taskApi";

type Props = {
  visible: boolean;
  task: Task | null;
  onClose: () => void;
  onSuccess: () => void;
};

export default function EditTaskModal({
  visible,
  task,
  onClose,
  onSuccess,
}: Props) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState<TaskStatus>("todo");
  const [assignee, setAssignee] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description);
      setStatus(task.status);
      setAssignee(task.assignee);
    }
  }, [task]);

  const handleSave = async () => {
    if (!task || saving) return;

    try {
      setSaving(true);

      await updateTask(task.id, {
        title,
        description,
        status,
        assignee,
      });

      onSuccess();
    } catch (error) {
      console.log(error);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent
      onRequestClose={onClose}
    >
      <View style={styles.overlay}>
        <View style={styles.container}>
          <Text style={styles.heading}>Edit Task</Text>

          <TextInput
            placeholder="Title"
            value={title}
            onChangeText={setTitle}
            style={styles.input}
          />

          <TextInput
            placeholder="Description"
            value={description}
            onChangeText={setDescription}
            style={styles.input}
            multiline
          />

          <TextInput
            placeholder="Assignee"
            value={assignee}
            onChangeText={setAssignee}
            style={styles.input}
          />

          <Picker
            selectedValue={status}
            onValueChange={(value) => setStatus(value as TaskStatus)}
          >
            <Picker.Item label="Todo" value="todo" />
            <Picker.Item label="In Progress" value="in_progress" />
            <Picker.Item label="Done" value="done" />
          </Picker>

          <View style={styles.button}>
            <Button title={saving ? "Saving..." : "Save"} onPress={handleSave} disabled={saving} />
          </View>

          <View style={styles.button}>
            <Button title="Cancel" color="red" onPress={onClose} />
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.4)",
    justifyContent: "center",
    alignItems: "center",
    padding: 20,
  },

  container: {
    width: "100%",
    backgroundColor: "#fff",
    borderRadius: 10,
    padding: 20,
  },

  heading: {
    fontSize: 20,
    fontWeight: "bold",
    marginBottom: 15,
    textAlign: "center",
  },

  input: {
    borderWidth: 1,
    borderColor: "#ccc",
    borderRadius: 8,
    padding: 10,
    marginBottom: 12,
  },

  button: {
    marginTop: 10,
  },
});
