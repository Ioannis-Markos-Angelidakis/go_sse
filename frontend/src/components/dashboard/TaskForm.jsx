import { useState } from 'react';
import {
    Card, CardContent,
    TextField, Button,
    FormControlLabel, Checkbox,
    Typography
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';

export default function TaskForm({ onCreateTask }) {
    const [newTask, setNewTask] = useState({
        title: '',
        content: '',
        public: false
    });

    const handleCreate = () => {
        if (!newTask.title.trim()) return;
        onCreateTask(newTask);
        setNewTask({ title: '', content: '', public: false });
    };

    return (
        <Card
            sx={{
                borderRadius: 3,
                boxShadow: '0 4px 15px rgba(0, 0, 0, 0.05)',
                height: '100%'
            }}
        >
            <CardContent>
                <Typography variant="h6" sx={{ mb: 1.5, fontWeight: 600, fontSize: '1.1rem' }}>
                    Create New Task
                </Typography>

                <TextField
                    label="Title"
                    value={newTask.title}
                    onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
                    fullWidth
                    size="small"
                    sx={{ mb: 1.5 }}
                />

                <TextField
                    label="Notes"
                    value={newTask.content}
                    onChange={(e) => setNewTask({ ...newTask, content: e.target.value })}
                    multiline
                    minRows={2}
                    fullWidth
                    size="small"
                    sx={{ mb: 1.5 }}
                />

                <FormControlLabel
                    control={
                        <Checkbox
                            checked={newTask.public}
                            onChange={(e) => setNewTask({ ...newTask, public: e.target.checked })}
                            color="primary"
                            size="small"
                        />
                    }
                    label="Make Public"
                    sx={{ mb: 1.5 }}
                />

                <Button
                    variant="contained"
                    color="primary"
                    onClick={handleCreate}
                    fullWidth
                    startIcon={<AddIcon />}
                    size="small"
                    sx={{ py: 0.8 }}
                >
                    Add Task
                </Button>
            </CardContent>
        </Card>
    );
}