import { useState, useEffect } from 'react';
import List from '@mui/material/List';
import { getList, deleteItem } from '../services/api';
import Item from './Item';

function ItemList({ listId }) {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch items when component mounts or listId changes
  useEffect(() => {
    fetchItems();
  }, [listId]);

  const fetchItems = async () => {
    try {
      setLoading(true);
      const response = await getList(listId); // Changed from getItems to getList
      
      // Items are nested in the response
      const itemsData = response.data.Items || [];
      setItems(itemsData);
      setError(null);
    } catch (err) {
      setError('Failed to fetch items');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    try {
      await deleteItem(id);
      // Remove from UI
      setItems(items.filter(item => item.id !== id));
    } catch (err) {
      setError('Failed to delete item');
      console.error(err);
    }
  };

  const handleEdit = (item) => {
    // TODO: Open edit form/modal
    console.log('Edit item:', item);
  };

  // make create function

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <List>
      {items.length === 0 ? (
        <div>No items yet</div>
      ) : (
        items.map(item => (
          <Item 
            key={item.id} 
            item={item} 
            onDelete={handleDelete}
            onEdit={handleEdit}
          />
        ))
      )}
    </List>
  );
}

export default ItemList;