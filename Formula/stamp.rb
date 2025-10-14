# This is a template Homebrew formula for the stamp CLI tool
# Copy this file to your homebrew-tap repository at Formula/stamp.rb
# GoReleaser will automatically update this formula when you create a new release

class Stamp < Formula
  desc "A simple Go CLI tool for generating note filenames based on date/time"
  homepage "https://github.com/totocaster/stamp"
  version "0.1.0"
  license "MIT"

  # For macOS
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/totocaster/stamp/releases/download/v0.1.0/stamp_0.1.0_Darwin_arm64.tar.gz"
      sha256 "TO_BE_UPDATED_BY_GORELEASER"
    else
      url "https://github.com/totocaster/stamp/releases/download/v0.1.0/stamp_0.1.0_Darwin_x86_64.tar.gz"
      sha256 "TO_BE_UPDATED_BY_GORELEASER"
    end
  # For Linux
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/totocaster/stamp/releases/download/v0.1.0/stamp_0.1.0_Linux_arm64.tar.gz"
      sha256 "TO_BE_UPDATED_BY_GORELEASER"
    else
      url "https://github.com/totocaster/stamp/releases/download/v0.1.0/stamp_0.1.0_Linux_x86_64.tar.gz"
      sha256 "TO_BE_UPDATED_BY_GORELEASER"
    end
  end

  # Optional: Go is only needed if building from source
  depends_on "go" => :optional

  def install
    # Install the stamp binary
    bin.install "stamp"

    # Create the nid symlink
    bin.install_symlink "stamp" => "nid"

    # Optional: Install shell completions if they exist
    # bash_completion.install "completions/stamp.bash" if File.exist?("completions/stamp.bash")
    # fish_completion.install "completions/stamp.fish" if File.exist?("completions/stamp.fish")
    # zsh_completion.install "completions/_stamp" if File.exist?("completions/_stamp")
  end

  test do
    # Test that the binary runs and shows version
    assert_match /stamp version/, shell_output("#{bin}/stamp version")

    # Test daily note generation
    assert_match(/\d{4}-\d{2}-\d{2}/, shell_output("#{bin}/stamp daily"))

    # Test that nid symlink works
    assert_predicate bin/"nid", :exist?
    assert_match(/\d{4}-\d{2}-\d{2}/, shell_output("#{bin}/nid daily"))

    # Test fleeting note generation
    assert_match(/\d{4}-\d{2}-\d{2}-F\d{6}/, shell_output("#{bin}/stamp fleeting"))

    # Test project note generation (should start with P)
    assert_match(/^P\d{4}/, shell_output("#{bin}/stamp project"))
  end
end