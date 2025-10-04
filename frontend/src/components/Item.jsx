import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';

function Item({ item, onDelete, onEdit }) {
  // Format the date nicely
  const formatDate = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  return (
    <ListItem
      secondaryAction={
        <>
          <IconButton edge="end" aria-label="edit" onClick={() => onEdit(item)}>
            <EditIcon />
          </IconButton>
          <IconButton edge="end" aria-label="delete" onClick={() => onDelete(item.id)}>
            <DeleteIcon />
          </IconButton>
        </>
      }
    >
      <ListItemText
        primary={`${item.title} ${item.item_date ? '- ' + formatDate(item.item_date) : ''}`}
        secondary={item.content}
      />
    </ListItem>
  );
}

export default Item;