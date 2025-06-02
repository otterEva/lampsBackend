import { FormLabel } from '@chakra-ui/form-control'
import { Input } from '@chakra-ui/react'

type PasswordInputProps = {
  password: string
  changeFunc: (e: React.ChangeEvent<HTMLInputElement>, cred: string) => void
}

const PasswordInput = ({ password, changeFunc }: PasswordInputProps) => {
  return (
    <FormLabel width="100%">
      <Input
        placeholder="password"
        type="password"
        value={password}
        onChange={(e) => changeFunc(e, 'password')}
        required
        borderRadius="10px"
        borderColor="rgba(52, 72, 255, 1)"
      />
    </FormLabel>
  )
}

export default PasswordInput
