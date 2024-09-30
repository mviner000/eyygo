import { spawn } from 'child_process';
import * as dotenv from 'dotenv';

dotenv.config({
  path: process.env.NODE_ENV === 'development' ? './.env.development' : './.env'
});

interface EnvVars {
  PORT: string;
  CERT_FILE: string;
  KEY_FILE: string;
  ALLOWED_ORIGINS: string;
}

const getEnvVars = (): EnvVars => {
  const PORT: string = process.env.PORT || '3000';
  const CERT_FILE: string = process.env.CERT_FILE || '';
  const KEY_FILE: string = process.env.KEY_FILE || '';
  let ALLOWED_ORIGINS: string = process.env.ALLOWED_ORIGINS || 'http://localhost:3001,https://eyymi.site';

  if (process.env.NODE_ENV === 'development') {
    ALLOWED_ORIGINS = '*';
  }

  return { PORT, CERT_FILE, KEY_FILE, ALLOWED_ORIGINS };
};

const command: string = process.env.NODE_ENV === 'development' ? './my-fiber-app.exe' : './my-fiber-app';

const envVars: EnvVars = getEnvVars();
const app = spawn(command, [], {
  env: {
    ...process.env,
    ...envVars
  }
});

app.stdout.on('data', (data: Buffer) => {
  console.log(`stdout: ${data}`);
});

app.stderr.on('data', (data: Buffer) => {
  console.error(`stderr: ${data}`);
});

app.on('close', (code: number) => {
  console.log(`child process exited with code ${code}`);
});