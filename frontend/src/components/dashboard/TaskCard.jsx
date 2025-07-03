import {
    Card, CardContent,
    Checkbox, TextField, IconButton,
    FormControlLabel, Box
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';

export default function TaskCard({ task, onUpdate, onDelete, isOwner }) {
    return (
        <Card
            sx={{
                mb: 1,
                borderLeft: task.public ? '3px solid' : 'none',
                borderColor: 'primary.main',
                backgroundColor: task.public ? 'rgba(234, 244, 255, 0.4)' : 'rgba(245, 245, 247, 0.5)',
                transition: 'all 0.3s cubic-bezier(0.25, 0.1, 0.25, 1)',
                '&:hover': {
                    transform: 'translateY(-1px)',
                    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.05)',
                }
            }}
        >
            {task.public && !isOwner && (
                <Box sx={{
                    backgroundColor: 'rgba(10, 132, 255, 0.1)',
                    color: 'primary.main',
                    padding: '2px 8px',
                    borderRadius: '12px',
                    fontSize: '0.7rem',
                    fontWeight: 500,
                    mx: 1,
                    mt: 1,
                    display: 'inline-block'
                }}>
                    By: {task.user?.email || 'Unknown'}
                </Box>
            )}

            <CardContent sx={{ pt: 1, pb: '8px !important' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 0.5 }}>
                    {isOwner && (
                        <Checkbox
                            checked={task.completed}
                            onChange={(e) => onUpdate(task.id, { completed: e.target.checked })}
                            color="primary"
                            size="small"
                            sx={{ mr: 0.5 }}
                        />
                    )}

                    <TextField
                        value={task.title}
                        onChange={(e) => isOwner && onUpdate(task.id, { title: e.target.value })}
                        variant="standard"
                        fullWidth
                        size="small"
                        InputProps={{
                            readOnly: !isOwner,
                            disableUnderline: true,
                            sx: {
                                fontSize: '0.95rem',
                                fontWeight: 500,
                                textDecoration: task.completed ? 'line-through' : 'none',
                                color: task.completed ? 'text.secondary' : 'text.primary',
                                '& .MuiInput-input': {
                                    padding: '1px 3px',
                                    '&:focus': {
                                        backgroundColor: 'rgba(10, 132, 255, 0.05)',
                                        borderRadius: '4px',
                                    }
                                }
                            }
                        }}
                    />

                    {isOwner && (
                        <IconButton
                            onClick={() => onDelete(task.id)}
                            color="error"
                            size="small"
                            sx={{ ml: 0.5 }}
                        >
                            <DeleteIcon fontSize="small" />
                        </IconButton>
                    )}
                </Box>

                <TextField
                    value={task.content || ''}
                    onChange={(e) => isOwner && onUpdate(task.id, { content: e.target.value })}
                    variant="outlined"
                    fullWidth
                    multiline
                    minRows={1}
                    maxRows={2}
                    placeholder="Add notes..."
                    size="small"
                    sx={{ mt: 1 }}
                    InputProps={{
                        readOnly: !isOwner,
                        sx: {
                            '& .MuiInputBase-input': {
                                padding: '4px 6px',
                                fontSize: '0.85rem',
                                '&:focus': {
                                    backgroundColor: 'rgba(10, 132, 255, 0.05)',
                                    borderRadius: '4px',
                                }
                            }
                        }
                    }}
                />

                {isOwner && (
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={task.public}
                                onChange={(e) => onUpdate(task.id, { public: e.target.checked })}
                                color="primary"
                                size="small"
                            />
                        }
                        label="Public"
                        sx={{ mt: 0.5, ml: 0, fontSize: '0.8rem' }}
                    />
                )}
            </CardContent>
        </Card>
    );
}