import { Card, CardContent, Typography } from '@mui/material';
import TaskCard from './TaskCard';

export default function TaskSection({
    title,
    tasks,
    onUpdate,
    onDelete,
    emptyText,
    isPersonal = false,
    currentUserId
}) {
    return (
        <Card sx={{ height: '100%', borderRadius: 3 }}>
            <CardContent sx={{ p: 1.5 }}>
                <Typography variant="h6" sx={{ fontWeight: 600, mb: 1, fontSize: '1.1rem' }}>
                    {title}
                </Typography>
                <div style={{ maxHeight: 'calc(100vh - 200px)', overflowY: 'auto' }}>
                    {tasks.length > 0 ? (
                        tasks.map(task => (
                            <TaskCard
                                key={task.id}
                                task={task}
                                onUpdate={onUpdate}
                                onDelete={onDelete}
                                isOwner={isPersonal || task.userId === currentUserId}
                            />
                        ))
                    ) : (
                        <Typography variant="body2" color="textSecondary" sx={{ py: 2, textAlign: 'center' }}>
                            {emptyText}
                        </Typography>
                    )}
                </div>
            </CardContent>
        </Card>
    );
}