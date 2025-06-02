import LoginInput from '../components/LoginInput'

type LoginFormProps = {
  setAdmin: (admin: boolean) => void
}

const LoginForm = ({ setAdmin }: LoginFormProps) => {
  return <LoginInput setAdmin={setAdmin} />
}

export default LoginForm
