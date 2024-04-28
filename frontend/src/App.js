import './App.css';
import React, { useState, useEffect } from 'react';
import Menu from './Menu/Menu.jsx';

function App() {
  const [loading, setLoading] = useState(true)
  const [recipes, setRecipes] = useState([""])

  useEffect(() => {
    fetch('/api/recipes')
    .then(response => response.json())
    .then(data => setRecipes(data))
    .finally(() => setLoading(false))
    // setLoading(false)
    // setRecipes(["Tiramisu", "Pasta", "Pizza quattro formagi"].sort())
  }, [])
  
  if (loading) {
    return <p>Loading...</p>
  }

  return (
    <Menu recipes={recipes}/>
  );
}

export default App;
