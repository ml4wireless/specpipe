# Nano React App Default Javascript Template

The default template project for [nano-react-app](https://github.com/nano-react-app/nano-react-app).

- `npm start` — This will spawn a development server with a default port of `5173`.
- `npm run build` — This will output a production build in the `dist` directory.
- `npm run preview` — This will run the production build locally with a default port of `5173` (this will not work if you haven't generated the production build yet).

## Custom port

You can use the `-p` flag to specify a port for development. To do this, you can either run `npm start` with an additional flag:

```
npm start -- --port 3000
```

Or edit the `start` script directly:

```
vite --port 3000
```