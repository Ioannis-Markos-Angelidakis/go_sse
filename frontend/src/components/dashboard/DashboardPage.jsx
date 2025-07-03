import { useState, useEffect } from 'react';
import { Grid, Container, CircularProgress } from '@mui/material';
import axios from '../../api/axios';
import { useAuth } from '../context/AuthContext';
import TaskForm from './TaskForm';
import TaskSection from './TaskSection';
import Header from '../layout/Header';

export default function DashboardPage() {
    const { user, logout } = useAuth();
    const [tasks, setTasks] = useState([]);
    const [publicTasks, setPublicTasks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [sseConnection, setSseConnection] = useState(null);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const [myTasks, publicTasks] = await Promise.all([
                    axios.get('http://192.168.0.13:5000/tasks'),
                    axios.get('http://192.168.0.13:5000/public-tasks'),
                ]);

                const enrichedPublicTasks = publicTasks.data.map(task => ({
                    ...task,
                    user: {
                        id: task.user.id,
                        email: task.user.email
                    }
                }));

                setTasks(myTasks.data);
                setPublicTasks(enrichedPublicTasks);
            } catch (error) {
                console.error('Failed to fetch tasks:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchTasks();

        const source = new EventSource('http://192.168.0.13:5000/events', {
            withCredentials: true
        });

        source.onopen = () => console.log('SSE connection established');

        source.onmessage = (event) => {
            if (!event.data || event.data.startsWith(':')) return;
            try {
                handleTaskEvent(JSON.parse(event.data));
            } catch (error) {
                console.error('Error processing SSE event:', error, event.data);
            }
        };

        source.onerror = (error) => {
            console.error('SSE connection error:', error);
            source.close();
        };

        setSseConnection(source);

        return () => source?.close();
    }, []);

    const handleTaskEvent = (event) => {
        switch (event.type) {
            case 'create':
                if (event.data?.public) {
                    setPublicTasks(prev => prev.some(t => t.id === event.data.id) ?
                        prev : [...prev, event.data]);
                }
                break;
            case 'update':
                setTasks(prev => prev.map(t => t.id === event.data.id ? { ...t, ...event.data } : t));
                setPublicTasks(prev => {
                    if (!event.data.public) return prev.filter(t => t.id !== event.data.id);
                    const existingIndex = prev.findIndex(t => t.id === event.data.id);
                    if (existingIndex >= 0) return prev.map(t =>
                        t.id === event.data.id ? { ...t, ...event.data } : t);
                    return event.data.public ? [...prev, event.data] : prev;
                });
                break;
            case 'delete':
                setTasks(prev => prev.filter(t => t.id !== event.taskId));
                setPublicTasks(prev => prev.filter(t => t.id !== event.taskId));
                break;
            default:
                console.warn('Unknown event type:', event.type);
        }
    };

    const handleCreateTask = async (newTask) => {
        try {
            const response = await axios.post('http://192.168.0.13:5000/tasks', newTask);
            setTasks(prev => [...prev, response.data]);
        } catch (error) {
            console.error('Failed to create task:', error);
        }
    };

    const handleUpdateTask = async (id, updates) => {
        try {
            await axios.put(`http://192.168.0.13:5000/tasks/${id}`, updates);
        } catch (error) {
            console.error('Failed to update task:', error);
        }
    };

    const handleDeleteTask = async (id) => {
        try {
            await axios.delete(`http://192.168.0.13:5000/tasks/${id}`);
        } catch (error) {
            console.error('Failed to delete task:', error);
        }
    };
    if (loading) {
        return (
            <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                <CircularProgress size={60} />
            </Container>
        );
    }

    return (
        <>
            <Header user={user} onLogout={logout} />

            <Container maxWidth="lg" sx={{ py: 2 }}>

                <Grid container spacing={2}>

                    <Grid item xs={12} md={3} sx={{ maxHeight: 300, maxWidth: 400 }}>
                        <TaskForm onCreateTask={handleCreateTask} />
                    </Grid>

                    <Grid item xs={12} md={9} container spacing={2}>
                        <Grid item >
                            <TaskSection
                                title="My Tasks"
                                tasks={tasks}
                                onUpdate={handleUpdateTask}
                                onDelete={handleDeleteTask}
                                emptyText="No personal tasks yet"
                                isPersonal={true}
                            />
                        </Grid>

                        <Grid item>
                            <TaskSection
                                title="Public Tasks"
                                tasks={publicTasks}
                                onUpdate={handleUpdateTask}
                                onDelete={handleDeleteTask}
                                emptyText="No public tasks yet"
                                currentUserId={user?.id}
                            />
                        </Grid>
                    </Grid>
                </Grid>
            </Container>
        </>
    );
}